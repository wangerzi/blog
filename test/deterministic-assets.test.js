'use strict';

const assert = require('node:assert/strict');
const { mkdtempSync, mkdirSync, writeFileSync } = require('node:fs');
const { tmpdir } = require('node:os');
const { join } = require('node:path');
const test = require('node:test');

const {
  appendVersion,
  createAssetVersioner,
  createHelper,
  normalizeAssetPath,
  withExtension,
} = require('../scripts/deterministic-assets');
const {
  normalizeJsonContent,
  sortArticleTags,
  sortFeedCategories,
} = require('../scripts/deterministic-output');

function fixture() {
  const baseDir = mkdtempSync(join(tmpdir(), 'blog-assets-'));
  const sourceDir = join(baseDir, 'source');
  const themeDir = join(baseDir, 'theme');
  mkdirSync(join(sourceDir, 'js'), { recursive: true });
  mkdirSync(join(themeDir, 'source', 'css', '_partial'), { recursive: true });
  writeFileSync(join(baseDir, '_config.yml'), 'theme: nexmoe\n');
  writeFileSync(join(sourceDir, 'js', 'app.js'), 'console.log("site")\n');
  writeFileSync(join(themeDir, 'source', 'css', 'style.styl'), '@import "_partial/base"\n');
  writeFileSync(join(themeDir, 'source', 'css', '_partial', 'base.styl'), 'body\n  color #222\n');
  return { baseDir, sourceDir, themeDir };
}

test('derives stable versions from bytes rather than file timestamps', () => {
  const paths = fixture();
  const first = createAssetVersioner(paths);
  const firstVersion = first('js/app.js', '.js');
  const secondVersion = createAssetVersioner(paths)('js/app.js', '.js');
  assert.equal(firstVersion, secondVersion);

  writeFileSync(join(paths.sourceDir, 'js', 'app.js'), 'console.log("changed")\n');
  assert.notEqual(createAssetVersioner(paths)('js/app.js', '.js'), firstVersion);
});

test('includes imported Stylus siblings in a CSS version', () => {
  const paths = fixture();
  const firstVersion = createAssetVersioner(paths)('css/style.css', '.css');
  writeFileSync(
    join(paths.themeDir, 'source', 'css', '_partial', 'base.styl'),
    'body\n  color #333\n',
  );
  assert.notEqual(
    createAssetVersioner(paths)('css/style.css', '.css'),
    firstVersion,
  );
});

test('renders compatible helpers with a deterministic query parameter', () => {
  const context = { url_for: (path) => `/${path}` };
  const css = createHelper({
    extension: '.css',
    tag: 'link',
    attribute: 'rel="stylesheet" href',
    versionFor: () => 'abc123',
  });
  const js = createHelper({
    extension: '.js',
    tag: 'script',
    attribute: 'src',
    versionFor: () => 'def456',
  });

  assert.equal(
    css.call(context, 'css/style'),
    '<link rel="stylesheet" href="/css/style.css?v=abc123">',
  );
  assert.equal(
    js.call(context, ['js/app', 'js/search.js']),
    '<script async src="/js/app.js?v=def456"></script>\n' +
      '<script async src="/js/search.js?v=def456"></script>',
  );
});

test('normalizes local paths and preserves existing query strings and hashes', () => {
  assert.equal(normalizeAssetPath('/css/style', '.css'), 'css/style.css');
  assert.equal(withExtension('/custom?theme=dark#top', '.css'), 'custom.css?theme=dark#top');
  assert.equal(
    appendVersion('/custom.css?theme=dark#top', 'abc123'),
    '/custom.css?theme=dark&v=abc123#top',
  );
  assert.throws(() => normalizeAssetPath('../secret', '.js'));
});

test('sorts generated metadata and search content deterministically', () => {
  assert.equal(
    sortArticleTags(
      '<meta property="article:tag" content="php">\n' +
        '<meta property="article:tag" content="Linux">\n',
    ),
    '<meta property="article:tag" content="Linux">\n' +
      '<meta property="article:tag" content="php">\n',
  );

  const normalized = JSON.parse(
    normalizeJsonContent(
      JSON.stringify({
        pages: [{ path: 'z' }, { path: 'a' }],
        posts: [
          { path: 'b', date: '2026-01-01', tags: [{ name: 'z' }, { name: 'a' }] },
          { path: 'a', date: '2026-01-01' },
          { path: 'new', date: '2026-02-01' },
        ],
        tags: [{ name: 'z' }, { name: 'a' }],
      }),
    ),
  );
  assert.deepEqual(
    normalized.pages.map((page) => page.path),
    ['a', 'z'],
  );
  assert.deepEqual(
    normalized.posts.map((post) => post.path),
    ['new', 'a', 'b'],
  );
  assert.deepEqual(
    normalized.posts[2].tags.map((tag) => tag.name),
    ['a', 'z'],
  );
  assert.deepEqual(
    normalized.tags.map((tag) => tag.name),
    ['a', 'z'],
  );

  assert.equal(
    sortFeedCategories(
      '<entry>\n  <category term="z"/>\n  <category term="a"/>\n</entry>',
    ),
    '<entry>\n  <category term="a"/>\n  <category term="z"/>\n</entry>',
  );
});
