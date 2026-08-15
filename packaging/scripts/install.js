import https from "node:https";
import http from "node:http";
import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import { execSync } from "node:child_process";

const GITHUB_REPO = "richavery/donk-cli-main";
const VERSION = process.env.npm_package_version || "latest";
const GITHUB_TOKEN = process.env.GITHUB_TOKEN || process.env.nPm_CONFIG_GITHUB_TOKEN || "";
const HOME = os.homedir();
const INSTALL_DIR = process.env.npm_config_global
  ? path.join(execSync("npm root -g").toString().trim(), ".bin")
  : path.join(process.cwd(), "node_modules", ".bin");
const TARGET = path.join(INSTALL_DIR, "donk-cli");

function api(url) {
  return new Promise((resolve, reject) => {
    const headers = { "User-Agent": "donk-cli-npm-installer" };
    if (GITHUB_TOKEN) headers["Authorization"] = `Bearer ${GITHUB_TOKEN}`;
    https.get(url, { headers }, (res) => {
      let data = "";
      res.on("data", (chunk) => (data += chunk));
      res.on("end", () => {
        try {
          resolve(JSON.parse(data));
        } catch (err) {
          reject(err);
        }
      });
    }).on("error", reject);
  });
}

function getRelease() {
  return api(`https://api.github.com/repos/${GITHUB_REPO}/releases/latest`);
}

function download(url) {
  return new Promise((resolve, reject) => {
    const mod = url.startsWith("https") ? https : http;
    const file = fs.createWriteStream(TARGET);
    mod.get(url, (res) => {
      if (res.statusCode && res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        download(res.headers.location).then(resolve).catch(reject);
        res.resume();
        return;
      }
      res.pipe(file);
      file.on("finish", () => {
        file.close();
        fs.chmodSync(TARGET, 0o755);
        resolve();
      });
    }).on("error", (err) => {
      fs.rmSync(TARGET, { force: true });
      reject(err);
    });
  });
}

async function main() {
  console.log(`Installing donk-cli@${VERSION || "latest"}...`);
  const release = await getRelease();
  const platform = os.platform();
  const arch = os.arch();
  const assetName =
    platform === "win32"
      ? `donk-cli_${release.tag_name}_Windows_${arch === "x64" ? "amd64" : arch}.zip`
      : `donk-cli_${release.tag_name}_${platform === "darwin" ? "darwin" : "linux"}_${arch === "x64" ? "amd64" : arch}.zip`;

  const asset = release.assets.find((item) => item.name === assetName);
  if (!asset) {
    console.error(`No matching asset found: ${assetName}`);
    process.exit(1);
  }

  await download(asset.browser_download_url);
  console.log(`Installed donk-cli to ${TARGET}`);
}

main().catch((err) => {
  console.error(err.message);
  process.exit(1);
});
