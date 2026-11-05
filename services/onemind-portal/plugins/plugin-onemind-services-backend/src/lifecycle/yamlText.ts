import YAML from 'yaml';

/**
 * Text-level helpers for multi-document Kubernetes YAML. They never re-serialise a
 * file: a document is kept or dropped byte for byte, so a review diff shows only the
 * documents that were removed.
 */

export type YamlChunk = {
  text: string;
  /** Parsed document, undefined for comment-only chunks or YAML that does not parse. */
  doc?: any;
};

/** Split at lines that are exactly `---`. Joining the chunks gives the original text. */
export function splitChunks(text: string): YamlChunk[] {
  const parts = text.split(/^(?=---[ \t]*\r?$)/m).filter(part => part !== '');
  return parts.map(part => {
    const body = part.replace(/^---[ \t]*\r?\n?/, '');
    let doc: any;
    try {
      const parsed = YAML.parse(body);
      if (parsed && typeof parsed === 'object') doc = parsed;
    } catch {
      doc = undefined;
    }
    return { text: part, doc };
  });
}

export const joinChunks = (chunks: YamlChunk[]) =>
  chunks.map(chunk => chunk.text).join('');

export type DocRef = { kind: string; name: string; namespace?: string };

export function docRef(doc: any): DocRef | undefined {
  if (!doc || typeof doc.kind !== 'string' || !doc.metadata?.name)
    return undefined;
  return {
    kind: doc.kind,
    name: String(doc.metadata.name),
    namespace: doc.metadata.namespace,
  };
}

/**
 * Remove the documents for which `match` returns true. Returns the new text (or
 * undefined when nothing is left that is a document, meaning the whole file can be
 * deleted) and the references that were removed.
 */
export function removeDocuments(
  text: string,
  match: (ref: DocRef, doc: any) => boolean,
): { text: string | undefined; removed: DocRef[]; kept: DocRef[] } {
  const chunks = splitChunks(text);
  const removed: DocRef[] = [];
  const kept: DocRef[] = [];
  const out: YamlChunk[] = [];
  for (const chunk of chunks) {
    const ref = docRef(chunk.doc);
    if (ref && match(ref, chunk.doc)) {
      removed.push(ref);
      continue;
    }
    if (ref) kept.push(ref);
    out.push(chunk);
  }
  if (!kept.length) return { text: undefined, removed, kept };
  return { text: joinChunks(out), removed, kept };
}

/** Every Kubernetes-style document in a file, with its reference. */
export function listDocuments(text: string): { ref: DocRef; doc: any }[] {
  return splitChunks(text)
    .map(chunk => ({ ref: docRef(chunk.doc), doc: chunk.doc }))
    .filter((x): x is { ref: DocRef; doc: any } => Boolean(x.ref));
}
