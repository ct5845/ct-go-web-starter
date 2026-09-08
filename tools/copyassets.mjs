// Copies the project's static files and the npm-installed frontend libraries into
// tmp/static, alongside the stylesheet Tailwind has already written there. npm
// runs this from the project root, so every path here is relative to it.
import { mkdirSync, readFileSync, readdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const outputDir = "tmp/static";
const staticDir = "internal/web/static";

// The single place that says where each npm-installed asset lands. Versions are
// pinned in package.json; adding a dependency means adding one line here.
const dependencies = [
  ["node_modules/alpinejs/dist/cdn.min.js", "alpine.min.js"],
  ["node_modules/htmx.org/dist/htmx.min.js", "htmx.min.js"],
];

// Neither fs.cpSync nor fs.copyFileSync is usable here: both copy the source's
// mode onto an existing destination, and chmod requires ownership rather than
// write permission. In a dev container the workspace is bind-mounted and
// tmp/static is owned by root, so those fail with EPERM on a directory the user
// can otherwise write to freely. Writing the bytes never touches metadata.
function copyDir(source, destination) {
  mkdirSync(destination, { recursive: true });

  for (const entry of readdirSync(source, { withFileTypes: true })) {
    const from = join(source, entry.name);
    const to = join(destination, entry.name);

    if (entry.isDirectory()) {
      copyDir(from, to);
    } else {
      writeFileSync(to, readFileSync(from));
    }
  }
}

copyDir(staticDir, outputDir);

for (const [source, name] of dependencies) {
  writeFileSync(join(outputDir, name), readFileSync(source));
}

console.log(
  `Copied ${staticDir} and ${dependencies.length} frontend dependencies to ${outputDir}`,
);
