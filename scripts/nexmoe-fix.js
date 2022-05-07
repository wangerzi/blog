hexo.extend.filter.register('after_render:html', (str, data) => {
    const cdn = data.config.theme_config.cdn;

    const replaceKeys = ['justifiedGallery', 'mdui', 'jquery', 'fancybox']
    for (let i = 0; i < replaceKeys.length; i++) {
        let key = replaceKeys[i];
        if (!cdn[key]) {
            continue;
        }
        if (cdn[key].js) {
            str = str.replace(new RegExp(`script src=".*cdn.jsdelivr.net.*${key}\.(umd|min|js).*?"`, 'gm'), `script src="${cdn[key].js}"`)
        }
        if (cdn[key].css) {
            str = str.replace(new RegExp(`href=".*cdn.jsdelivr.net.*${key}.*"`, 'gm'), `href="${cdn[key].css}"`)
        }
    }
    return str;
})