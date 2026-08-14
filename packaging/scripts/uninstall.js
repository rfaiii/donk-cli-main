import fs from "node:fs";
import path from "node:path";

const TARGET = process.argv[1] || path.join(process.cwd(), "node_modules", ".bin", "donk-cli");
fs.rmSync(TARGET, { force: true });
console.log("Removed donk-cli");
