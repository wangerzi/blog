'use strict';

const { createHash } = require('crypto');
const { existsSync, readFileSync, readdirSync, statSync } = require('fs');
const { dirname, extname, join, relative, resolve, sep } = require('path');

function normalizeAssetPath(assetPath, extension) {
  const [withoutHash] = String(assetPath).split('#', 1);
  const [pathname] = withoutHash.split('?', 1);
  const withExtension = extname(pathname) ? pathname : `${pathname}${extension}`;
  const normalized = withExtension.replace(/^\/+/, '').replace(/\\/g, '/');

  if (
    !normalized ||
    normalized.split('/').includes('..') ||
    /^[a-z][a-z\d+.-]*:/i.test(normalized)
  ) {
    throw new Error(`Unsupported local asset path: ${assetPath}`);
  }

  return normalized;
}

function appendVersion(url, version) {
  const hashIndex = url.indexOf('#');
  const beforeHash = hashIndex === -1 ? url : url.slice(0, hashIndex);
  const hash = hashIndex === -1 ? '' : url.slice(hashIndex);
  const separator = beforeHash.includes('?') ? '&' : '?';
  return `${beforeHash}${separator}v=${version}${hash}`;
}

function withExtension(assetPath, extension) {
  const value = String(assetPath);
  const suffixIndex = value.search(/[?#]/);
  const pathname = suffixIndex === -1 ? value : value.slice(0, suffixIndex);
  const suffix = suffixIndex === -1 ? '' : value.slice(suffixIndex);
  const normalized = normalizeAssetPath(pathname, extension);
  return `${normalized}${suffix}`;
}

function addFile(hash, root, path) {
  const details = statSync(path);
  if (details.isDirectory()) {
    for (const entry of readdirSync(path).sort()) {
      addFile(hash, root, join(path, entry));
    }
    return;
  }
  if (!details.isFile()) return;

  hash.update(relative(root, path).split(sep).join('/'));
  hash.update('\0');
  hash.update(readFileSync(path));
  hash.update('\0');
}

function createAssetVersioner({ baseDir, sourceDir, themeDir }) {
  const configPaths = [
    join(baseDir, '_config.yml'),
    join(baseDir, '_config.nexmoe.yml'),
    join(themeDir, '_config.yml'),
  ].filter(existsSync);

  return (assetPath, extension) => {
    const normalized = normalizeAssetPath(assetPath, extension);
    const hash = createHash('sha256');
    hash.update(normalized);
    hash.update('\0');

    for (const configPath of configPaths) addFile(hash, baseDir, configPath);

    const siteAsset = resolve(sourceDir, normalized);
    const themeAsset = resolve(themeDir, 'source', normalized);
    if (siteAsset.startsWith(`${resolve(sourceDir)}${sep}`) && existsSync(siteAsset)) {
      addFile(hash, sourceDir, siteAsset);
    } else if (
      themeAsset.startsWith(`${resolve(themeDir, 'source')}${sep}`) &&
      existsSync(themeAsset)
    ) {
      addFile(hash, join(themeDir, 'source'), themeAsset);
    } else if (extension === '.css') {
      const stylusEntry = resolve(
        themeDir,
        'source',
        normalized.replace(/\.css$/, '.styl'),
      );
      if (existsSync(stylusEntry)) {
        // Stylus entries can import sibling partials, so hash the complete tree.
        addFile(hash, join(themeDir, 'source'), dirname(stylusEntry));
      }
    }

    return hash.digest('hex').slice(0, 16);
  };
}

function createHelper({ extension, tag, attribute, versionFor }) {
  function helper(...assets) {
    return assets.reduce((result, asset, index) => {
      const separator = index ? '\n' : '';
      if (Array.isArray(asset)) {
        return result + separator + Reflect.apply(helper, this, asset);
      }

      const normalized = normalizeAssetPath(asset, extension);
      const url = appendVersion(
        this.url_for(withExtension(asset, extension)),
        versionFor(normalized, extension),
      );
      const asyncAttribute = tag === 'script' ? ' async' : '';
      const close = tag === 'script' ? `></${tag}>` : '>';
      return `${result}${separator}<${tag}${asyncAttribute} ${attribute}="${url}"${close}`;
    }, '');
  }

  return helper;
}

function registerDeterministicAssetHelpers(hexoInstance) {
  hexoInstance.extend.filter.register('before_generate', () => {
    const versionFor = createAssetVersioner({
      baseDir: hexoInstance.base_dir,
      sourceDir: hexoInstance.source_dir,
      themeDir: hexoInstance.theme_dir,
    });

    hexoInstance.extend.helper.register(
      'css_auto_version',
      createHelper({
        extension: '.css',
        tag: 'link',
        attribute: 'rel="stylesheet" href',
        versionFor,
      }),
    );
    hexoInstance.extend.helper.register(
      'js_auto_version',
      createHelper({
        extension: '.js',
        tag: 'script',
        attribute: 'src',
        versionFor,
      }),
    );
  });
}

if (typeof hexo !== 'undefined') registerDeterministicAssetHelpers(hexo);

module.exports = {
  appendVersion,
  createAssetVersioner,
  createHelper,
  normalizeAssetPath,
  registerDeterministicAssetHelpers,
  withExtension,
};
