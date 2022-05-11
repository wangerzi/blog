import imagemin from 'imagemin';
import imageminJpegtran from 'imagemin-jpegtran';
import imageminPngquant from 'imagemin-pngquant';

async function handle() {
    return  await imagemin(['source/static/assets/*.{jpg,png}'], {
        destination: 'source/static/assets/',
        plugins: [
            imageminJpegtran(),
            imageminPngquant({
                quality: [0.6, 0.8]
            })
        ]
    });
}

handle().then(files => {
    console.log('handled', files.map(v => v.destinationPath));
})