const COS = require("cos-nodejs-sdk-v5");
const fs = require("fs");
const path = require("path");
const cos = new COS({
  SecretId: process.env.TENCENT_SECRET_ID,
  SecretKey: process.env.TENCENT_SECRET_KEY,
});

const uploadDir = (dir, bucket, region, prefix) => {
  const files = fs.readdirSync(dir).map(file => ({
    Bucket: bucket,
    Region: region,
    Key: `${prefix}/${file}`,
    FilePath: path.join(dir, file)  // 使用 FilePath 以便使用 uploadFiles 方法
  }));

  // 使用 uploadFiles 方法批量上传
  cos.uploadFiles({
    files: files,
    SliceSize: 1024 * 1024,  // 1MB
    onProgress: (progressData) => {
      console.log(JSON.stringify(progressData));
    },
    async: true,
  }, (err, data) => {
    if (err) {
      console.error(err);
    } else {
      console.log("上传完成:", data);
    }
  });
};

uploadDir(
  "public",
  process.env.COS_BUCKET,
  process.env.COS_REGION,
  process.env.COS_UPLOAD_PATH
);
