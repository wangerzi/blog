'use strict';

function compareText(left, right) {
  const a = String(left ?? '');
  const b = String(right ?? '');
  return a < b ? -1 : a > b ? 1 : 0;
}

function entityKey(value) {
  if (typeof value === 'string') return value;
  if (!value || typeof value !== 'object') return '';
  return value.path || value.slug || value.name || value.title || '';
}

function sortTaxonomy(value) {
  if (!Array.isArray(value)) return value;
  return [...value].sort((left, right) =>
    compareText(entityKey(left), entityKey(right)),
  );
}

function normalizeItem(item) {
  if (!item || typeof item !== 'object') return item;
  const result = { ...item };
  if (Array.isArray(result.tags)) result.tags = sortTaxonomy(result.tags);
  if (Array.isArray(result.categories)) {
    result.categories = sortTaxonomy(result.categories);
  }
  return result;
}

function normalizeJsonContent(serialized) {
  const content = JSON.parse(serialized);
  if (Array.isArray(content)) {
    return JSON.stringify(
      content.map(normalizeItem).sort((left, right) =>
        compareText(entityKey(left), entityKey(right)),
      ),
    );
  }

  if (Array.isArray(content.pages)) {
    content.pages = content.pages
      .map(normalizeItem)
      .sort((left, right) => compareText(left.path, right.path));
  }
  if (Array.isArray(content.posts)) {
    content.posts = content.posts.map(normalizeItem).sort((left, right) => {
      const byDate = compareText(right.date, left.date);
      return byDate || compareText(left.path, right.path);
    });
  }
  if (Array.isArray(content.tags)) content.tags = sortTaxonomy(content.tags);
  if (Array.isArray(content.categories)) {
    content.categories = sortTaxonomy(content.categories);
  }

  return JSON.stringify(content);
}

function sortArticleTags(html) {
  return html.replace(
    /(?:<meta property="article:tag" content="[^"]*">\n)+/g,
    (block) => `${block.trimEnd().split('\n').sort(compareText).join('\n')}\n`,
  );
}

function sortFeedCategories(xml) {
  return xml.replace(/<entry>[\s\S]*?<\/entry>/g, (entry) => {
    const pattern = /^[ \t]*<category\b[^\n]*\/>[ \t]*$/gm;
    const categories = [...entry.matchAll(pattern)].map((match) => match[0]);
    if (categories.length < 2) return entry;

    let index = 0;
    categories.sort(compareText);
    return entry.replace(pattern, () => categories[index++]);
  });
}

function normalizeGeneratorResult(result) {
  if (Array.isArray(result)) return result.map(normalizeGeneratorResult);
  if (!result || typeof result !== 'object' || typeof result.data !== 'string') {
    return result;
  }
  if (result.path === 'content.json') {
    return { ...result, data: normalizeJsonContent(result.data) };
  }
  if (String(result.path).endsWith('.xml')) {
    return { ...result, data: sortFeedCategories(result.data) };
  }
  return result;
}

function registerDeterministicOutput(hexoInstance) {
  hexoInstance.extend.filter.register('after_render:html', sortArticleTags);

  let generatorsWrapped = false;
  hexoInstance.extend.filter.register('before_generate', () => {
    if (generatorsWrapped) return;
    const feedTypes = Array.isArray(hexoInstance.config.feed?.type)
      ? hexoInstance.config.feed.type
      : [hexoInstance.config.feed?.type];
    for (const name of ['json-content', ...feedTypes].filter(Boolean)) {
      const original = hexoInstance.extend.generator.get(name);
      if (!original) continue;
      hexoInstance.extend.generator.register(name, async function (site) {
        const result = await Reflect.apply(original, this, [site]);
        return normalizeGeneratorResult(result);
      });
    }
    generatorsWrapped = true;
  });
}

if (typeof hexo !== 'undefined') registerDeterministicOutput(hexo);

module.exports = {
  normalizeGeneratorResult,
  normalizeJsonContent,
  registerDeterministicOutput,
  sortArticleTags,
  sortFeedCategories,
};
