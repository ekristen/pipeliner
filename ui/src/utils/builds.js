const collator = new Intl.Collator("en", {
  numeric: true,
  sensitivity: "base",
});

export const sort = (a, b) => {
  if (a.stage_idx < b.stage_idx) return -1;
  if (a.stage_idx > b.stage_idx) return 1;
  if (collator.compare(a.name, b.name) < 0) return -1;
  if (collator.compare(a.name, b.name) > 0) return 1;
  if (a.id < b.id) return 1;
  if (a.id > b.id) return -1;
  return 0;
}