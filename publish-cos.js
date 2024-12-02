const COS = require("cos-nodejs-sdk-v5");
const fs = require("fs");
const path = require("path");
const cos = new COS({
  SecretId: process.env.TENCENT_SECRET_ID,
  SecretKey: process.env.TENCENT_SECRET_KEY,
});
const uploadDir = (dir, bucket, region, prefix) => {
  fs.readdirSync(dir).forEach((file) => {
    const filePath = path.join(dir, file);
    cos.putObject(
      {
        Bucket: bucket,
        Region: region,
        Key: `${prefix}/${file}`,
        Body: fs.readFileSync(filePath),
      },
      (err, data) => {
        if (err) console.error(err);
        else console.log(`Uploaded: ${file}`);
      }
    );
  });
};
uploadDir(
  "public",
  process.env.COS_BUCKET,
  process.env.COS_REGION,
  process.env.COS_UPLOAD_PATH
);
