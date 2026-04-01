export type TableEntry = {
  name: string;
  type: string; // "table" or "view"
};

export type ColumnEntry = {
  name: string;
  type: string;
  pk: boolean;
  notNull: boolean;
  comment: string;
};
