const path = require('path')
const cheerio = require('cheerio');
const imageRoot = hexo.config.imageRoot || '/';
const sourcePath = hexo.config.source_dir;

function convertRoot(docAbsPath, relSrc, sourcePath='') {
    const dirname = path.dirname(docAbsPath)
    const srcPath = path.normalize(path.resolve(dirname, relSrc));

    return path.relative(sourcePath, srcPath).replace(/\\/g, '/')
}

function urlForHelper(path = '/', docPath = '') {
    path = path.replace('\\', '/');
    if (/http[s]*.*|\/\/.*/.test(path) ||
        /^\s+\//.test(path)) {
        return path
    }
    if (path[0] === '#' || path.startsWith('//')) {
        return path;
    }

    if (path.startsWith('/')) {
        path = path.slice(1)
    }

    path = convertRoot(docPath, path, sourcePath)
    // Prepend path
    path = imageRoot + path;

    // path.replace(/\/{2,}/g, '/');
    // return path.replace(/(\\|\/){2,}/g, '/');
    return path;
}
hexo.extend.filter.register('after_post_render', function (data) {
    const docPath = data.full_source
    if (!docPath.endsWith('.md')) {
        console.log('path break', data.asset_dir, data.full_source)
        return ;
    }
    // console.log('path', data.asset_dir, data.full_source)
    // console.log("excerpt === " + data.excerpt);
    // console.log("more === " + data.more);
    // console.log("content === " + data.content);
    if (data.cover) {
        data.cover = urlForHelper(data.cover, docPath)
    }
    const dispose = ['excerpt', 'more', 'content'];
    for (var i = 0; i < dispose.length; i++) {
        var key = dispose[i];

        if (!data[key]) {
            continue;
        }
        var $ = cheerio.load(data[key], {
            ignoreWhitespace: false,
            xmlMode: false,
            lowerCaseTags: false,
            decodeEntities: false
        });

        $('img').each(function () {
            var src = $(this).attr('src');
            if (src) {
                // For windows style path, we replace '\' to '/'.
                const newUrl = urlForHelper(src, docPath);
                // console.log('change', $(this).parent().attr('href'), src, newUrl, data.full_source);
                if (src !== newUrl) {
                    $(this).attr('src', newUrl);
                }
                // too slow!!!
                if (src.indexOf($(this).parent().attr('href')) !== -1) {
                    $(this).parent().attr('href', newUrl)
                }
            } else {
            }
        });
        data[key] = $.html();
    }
    // if (this.config.imagePrefix) {
    //     data.content = data.content.replace(
    //         new RegExp(themeCfg.imageCDN.origin, "gm"),
    //         themeCfg.imageCDN.to
    //     );
    // }
})

console.log('')