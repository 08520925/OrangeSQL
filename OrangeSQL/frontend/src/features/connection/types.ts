export type DriverType = 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';

export type ConnectionProfile = {
  id: string;
  name: string;
  driver: DriverType;
  createdAt: string;

  // SQLite
  path?: string;

  // TCP 系共通
  host?: string;
  port?: number;
  user?: string;
  password?: string;
  dbName?: string;
  sslMode?: string;
};

export type CreateProfileRequest = {
  name: string;
  driver: DriverType;
  path?: string;
  host?: string;
  port?: number;
  user?: string;
  password?: string;
  dbName?: string;
  sslMode?: string;
};

export type UpdateProfileRequest = {
  name?: string;
  path?: string;
  host?: string;
  port?: number;
  user?: string;
  password?: string;
  dbName?: string;
  sslMode?: string;
};

export const DRIVERS = [
  { value: 'sqlite' as const, label: 'SQLite' },
  { value: 'postgres' as const, label: 'PostgreSQL' },
  { value: 'mysql' as const, label: 'MySQL' },
  { value: 'sqlserver' as const, label: 'SQL Server' },
];

export const DEFAULT_PORTS: Record<DriverType, number> = {
  sqlite: 0,
  postgres: 5432,
  mysql: 3306,
  sqlserver: 1433,
};
