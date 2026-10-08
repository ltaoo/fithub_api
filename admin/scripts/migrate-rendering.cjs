// Mechanical import migration; does not change page/business implementation.
const fs = require("node:fs");
const path = require("node:path");
const ts = require("typescript");
const used_icons = new Set();
function walk(directory) {
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const filename = path.join(directory, entry.name);
    if (entry.isDirectory()) { if (filename !== "src/timeless") walk(filename); continue; }
    if (!/\.(ts|tsx)$/.test(filename)) continue;
    const source = fs.readFileSync(filename, "utf8");
    for (const match of source.matchAll(/import\s*\{([^}]+)\}\s*from\s*["'](?:lucide-solid|@\/timeless\/icons)["']/g)) {
      match[1].split(",").map(name => name.trim().split(/\s+as\s+/)[0]).filter(Boolean).forEach(name => used_icons.add(name));
    }
    const next = source.replace(/(["'])solid-js(?:\/web|\/jsx-runtime)?\1/g, '"@/timeless"')
      .replace(/(["'])lucide-solid\1/g, '"@/timeless/icons"');
    if (next !== source) fs.writeFileSync(filename, next);
  }
}
walk("src");
const source = fs.readFileSync("node_modules/lucide-solid/dist/esm/lucide-solid.js", "utf8");
const ast = ts.createSourceFile("icons.js", source, ts.ScriptTarget.ESNext, true, ts.ScriptKind.JS);
const declarations = new Map();
for (const statement of ast.statements) {
  if (ts.isVariableStatement(statement)) for (const declaration of statement.declarationList.declarations)
    declarations.set(declaration.name.getText(ast), declaration.initializer);
}
const exports_map = new Map();
for (const statement of ast.statements) if (ts.isExportDeclaration(statement) && statement.exportClause && ts.isNamedExports(statement.exportClause)) {
  for (const element of statement.exportClause.elements) exports_map.set(element.name.text, element.propertyName?.text || element.name.text);
}
function literal(node) {
  if (ts.isStringLiteral(node) || ts.isNumericLiteral(node)) return ts.isNumericLiteral(node) ? Number(node.text) : node.text;
  if (ts.isArrayLiteralExpression(node)) return node.elements.map(literal);
  if (ts.isObjectLiteralExpression(node)) return Object.fromEntries(node.properties.map(property => [property.name.text, literal(property.initializer)]));
  throw new Error("Unexpected icon data " + node.getText(ast));
}
let output = '// Lucide 0.314.0 SVG data (ISC); rendered by Timeless, no Solid runtime.\nimport { h } from "./index";\n';
for (const name of [...used_icons].sort()) {
  let node = declarations.get(exports_map.get(name) || name);
  if (node && ts.isIdentifier(node)) node = declarations.get(node.text);
  const match = node?.getText(ast).match(/iconNode:\s*([\w$]+)/);
  if (!match) throw new Error("Cannot migrate icon " + name);
  const data = literal(declarations.get(match[1]));
  output += `export const ${name} = (props: any = {}) => h("svg", { viewBox: "0 0 24 24", width: 24, height: 24, fill: "none", stroke: "currentColor", "stroke-width": 2, "stroke-linecap": "round", "stroke-linejoin": "round", ...props, children: ${JSON.stringify(data)}.map(([tag, attributes]: any) => h(tag, attributes)) });\n`;
}
fs.writeFileSync("src/timeless/icons.ts", output);
console.log(`Migrated rendering imports and ${used_icons.size} icons.`);
