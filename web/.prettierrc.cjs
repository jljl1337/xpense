// @ts-check

/** @type {import("prettier").Config} */
module.exports = {
  // Since prettier 3.0, manually specifying plugins is required
  plugins: ["@ianvs/prettier-plugin-sort-imports"],
  // This plugin's options
  importOrder: [
    "^react$",
    "",
    "<THIRD_PARTY_MODULES>",
    "",
    "^~/components/ui/(.*)$",
    "^~/lib/utils$",
    "",
    "^~/(.*)$",
  ],
  importOrderParserPlugins: ["typescript", "jsx", "decorators-legacy"],
  importOrderTypeScriptVersion: "5.0.0",
  importOrderCaseSensitive: false,
};