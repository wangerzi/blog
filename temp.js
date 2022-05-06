const path = require('path')

function convertRoot(docAbsPath, relSrc, sourcePath='') {
    const dirname = path.dirname(docAbsPath)
    const srcPath = path.normalize(path.resolve(dirname, relSrc));

    return path.relative(sourcePath, srcPath).replace(/\\/g, '/')
}

const res = convertRoot('D:\\phpStudy\\WWW\\github\\blog\\source\\_posts\\alchemy、geth-推荐的-gasprice-是怎么来的.md', '../static/uploads/git.jpg', 'source')
console.log(res);