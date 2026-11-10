// Curated reference/planning cards are reviewed by hand and must survive every inventory refresh unchanged.
// Rule: a Component in apps.yaml tagged `reference` or `planned` is curated. Discovery, enrichment and
// verification may rewrite everything else, but these cards are put back exactly as they were
// (raw YAML text, so folded descriptions and comments are kept too). The catalog tests (audit, topology,
// tak-ecosystem, ...) pin these cards, which is how the 2026-11-08 refresh broke the image build.
const fs = require('fs');
const cp = require('child_process');
const YAML = require('yaml');

const CURATED_TAGS = ['reference', 'planned'];
// Document separator lines only; a `---` inside a block scalar is indented, so it never matches.
const split = text => text.split(/^---\n/m);
const nameOf = block => {
  try { return YAML.parse(block)?.metadata?.name; } catch { return undefined; }
};
const isCuratedBlock = block => {
  try {
    const entity = YAML.parse(block);
    return entity?.kind === 'Component' && (entity.metadata?.tags || []).some(tag => CURATED_TAGS.includes(tag));
  } catch { return false; }
};

// name -> raw YAML block text of every curated card in `text`.
function curatedBlocks(text) {
  const blocks = new Map();
  for (const block of split(text)) if (isCuratedBlock(block)) blocks.set(nameOf(block), block);
  return blocks;
}

// Rewrite `refreshed` so each curated card from `original` is back byte-for-byte. A curated card the refresh
// dropped is appended again. Everything else in `refreshed` is untouched.
function restoreCuratedCards(original, refreshed) {
  const saved = curatedBlocks(original);
  const seen = new Set();
  const blocks = split(refreshed).map(block => {
    const name = nameOf(block);
    if (saved.has(name) && isCuratedBlock(saved.get(name))) { seen.add(name); return saved.get(name); }
    return block;
  });
  for (const [name, block] of saved) if (!seen.has(name)) blocks.push(block);
  return blocks.map(block => (block.endsWith('\n') ? block : block + '\n')).join('---\n');
}

// Call at the end of a script that rewrote apps.yaml. `originals` maps file path -> text read before the run.
function preserveCuratedCards(originals) {
  for (const [file, original] of originals) {
    if (original === null || original === undefined || !fs.existsSync(file)) continue;
    const refreshed = fs.readFileSync(file, 'utf8');
    const restored = restoreCuratedCards(original.toString(), refreshed);
    if (restored !== refreshed) fs.writeFileSync(file, restored);
  }
}

// The exact catalog test command the image build runs (packages/backend/Dockerfile), so a bad refresh fails
// here and not in build-onemind-portal.yml.
function dockerfileTestArgs(dockerfile = 'packages/backend/Dockerfile') {
  const line = fs.readFileSync(dockerfile, 'utf8').split('\n').find(l => /^RUN node --test /.test(l));
  if (!line) throw new Error(`No "RUN node --test" step found in ${dockerfile}`);
  return line.replace(/^RUN node /, '').trim().split(/\s+/);
}
function runDockerfileCatalogTests(dockerfile) {
  cp.execFileSync(process.execPath, dockerfileTestArgs(dockerfile), { stdio: 'inherit', timeout: 300000 });
}

module.exports = { curatedBlocks, restoreCuratedCards, preserveCuratedCards, dockerfileTestArgs, runDockerfileCatalogTests };
