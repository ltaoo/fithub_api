import ts from "typescript";
import type { Plugin } from "vite";

// Preserve lazy JSX properties while targeting Timeless VNodes, not Solid DOM helpers.
export function timeless_jsx(): Plugin {
  return {
    name: "admin-timeless-jsx",
    enforce: "pre",
    transform(code, id) {
      if (!id.split("?")[0].endsWith(".tsx")) return;
      const factory = ts.factory;
      const getter = (name: string, expression: ts.Expression) => factory.createGetAccessorDeclaration(
        undefined, factory.createStringLiteral(name), [], undefined,
        factory.createBlock([factory.createReturnStatement(expression)], true),
      );
      const arrow = (expression: ts.Expression) => factory.createArrowFunction(undefined, undefined, [], undefined,
        factory.createToken(ts.SyntaxKind.EqualsGreaterThanToken), expression);
      const transformer: ts.TransformerFactory<ts.SourceFile> = context => {
        const visit: ts.Visitor = node => {
          if (ts.isJsxElement(node) || ts.isJsxSelfClosingElement(node) || ts.isJsxFragment(node)) {
            const fragment = ts.isJsxFragment(node);
            const opening = fragment ? undefined : ts.isJsxElement(node) ? node.openingElement : node;
            const tag = opening?.tagName;
            const tag_text = tag?.getText();
            const component = fragment ? factory.createIdentifier("__tt_fragment")
              : tag_text && /^[a-z]/.test(tag_text) && !tag_text.includes(".")
                ? factory.createStringLiteral(tag_text)
                : tag as ts.Expression;
            const properties: ts.ObjectLiteralElementLike[] = [];
            for (const attr of opening?.attributes.properties || []) {
              if (ts.isJsxSpreadAttribute(attr)) {
                properties.push(factory.createSpreadAssignment(ts.visitNode(attr.expression, visit) as ts.Expression));
                continue;
              }
              const name = attr.name.getText();
              let expression: ts.Expression = factory.createTrue();
              if (attr.initializer && ts.isStringLiteral(attr.initializer)) expression = attr.initializer;
              if (attr.initializer && ts.isJsxExpression(attr.initializer) && attr.initializer.expression) {
                expression = ts.visitNode(attr.initializer.expression, visit) as ts.Expression;
              }
              if (name === "ref" && ts.isIdentifier(expression)) {
                expression = factory.createArrowFunction(undefined, undefined,
                  [factory.createParameterDeclaration(undefined, undefined, "__node")], undefined,
                  factory.createToken(ts.SyntaxKind.EqualsGreaterThanToken),
                  factory.createAssignment(expression, factory.createIdentifier("__node")));
              }
              properties.push(getter(name, expression));
            }
            const child_nodes = (ts.isJsxSelfClosingElement(node) ? [] : [...node.children]).filter(child =>
              !(ts.isJsxText(child) && !child.text.trim()) && !(ts.isJsxExpression(child) && !child.expression));
            const children: ts.Expression[] = [];
            for (const child of child_nodes) {
              if (ts.isJsxText(child)) {
                const lines = child.text.replace(/\r/g, "").split("\n");
                const text = lines.map((line, index) => index === 0 ? line.trimEnd() : index === lines.length - 1 ? line.trimStart() : line.trim()).filter(Boolean).join(" ");
                if (text) children.push(factory.createStringLiteral(text));
              } else if (ts.isJsxExpression(child)) {
                if (child.expression) children.push(factory.createCallExpression(factory.createIdentifier("__tt_dynamic"), undefined,
                  [arrow(ts.visitNode(child.expression, visit) as ts.Expression)]));
              } else {
                children.push(factory.createCallExpression(factory.createIdentifier("__tt_lazy"), undefined,
                  [arrow(ts.visitNode(child, visit) as ts.Expression)]));
              }
            }
            if (children.length) {
              // Control-flow render functions are passed directly, not mounted as children.
              const single = child_nodes.length === 1 && ts.isJsxExpression(child_nodes[0]) && child_nodes[0].expression;
              const value = single && ts.isArrowFunction(single)
                ? ts.visitNode(single, visit) as ts.Expression
                : factory.createArrayLiteralExpression(children);
              properties.push(getter("children", value));
            }
            return factory.createCallExpression(factory.createIdentifier("__tt_h"), undefined,
              [component, factory.createObjectLiteralExpression(properties, true)]);
          }
          return ts.visitEachChild(node, visit, context);
        };
        return source => ts.visitNode(source, visit) as ts.SourceFile;
      };
      const result = ts.transpileModule(code, {
        fileName: id,
        compilerOptions: { target: ts.ScriptTarget.ESNext, module: ts.ModuleKind.ESNext, sourceMap: true, jsx: ts.JsxEmit.Preserve },
        transformers: { before: [transformer] },
      });
      return {
        code: 'import { h as __tt_h, Fragment as __tt_fragment, dynamic as __tt_dynamic, lazy as __tt_lazy } from "@/timeless";\n' + result.outputText,
        map: null,
      };
    },
  };
}
