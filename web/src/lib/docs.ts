import { promises as fs } from "fs";
import path from "path";

/**
 * Load markdown content from the docs content directory
 * @param pathSegments - Path segments relative to /src/content/docs/
 * @returns The markdown content as a string
 */
export async function getDocsContent(...pathSegments: string[]): Promise<string> {
  const filePath = path.join(
    process.cwd(),
    "src/content/docs",
    ...pathSegments
  );
  
  const content = await fs.readFile(filePath, "utf8");
  return content;
}

