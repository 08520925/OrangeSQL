import { ref, onMounted } from "vue";
import {
  fetchProfiles,
  createProfile,
  updateProfile,
  deleteProfile,
  connectProfile,
} from "../../shared/api";
import type { ConnectionProfile } from "./types";

/**
 * 接続プロファイルの管理を担う composable。
 */
export function useConnection() {
  const profiles = ref<ConnectionProfile[]>([]);
  const activeId = ref<string>("");
  const loading = ref<boolean>(false);

  async function refresh(): Promise<void> {
    loading.value = true;
    try {
      const res = await fetchProfiles();
      profiles.value = res.profiles;
      activeId.value = res.activeId;
    } catch {
      // エラー時は何もしない
    } finally {
      loading.value = false;
    }
  }

  async function create(name: string, driver: string, path: string): Promise<void> {
    await createProfile({ name, driver, path });
    await refresh();
  }

  async function update(id: string, name: string, path: string): Promise<void> {
    await updateProfile(id, { name, path });
    await refresh();
  }

  async function remove(id: string): Promise<void> {
    await deleteProfile(id);
    await refresh();
  }

  async function connect(id: string): Promise<void> {
    await connectProfile(id);
    activeId.value = id;
    await refresh();
  }

  onMounted(() => {
    void refresh();
  });

  return { profiles, activeId, loading, refresh, create, update, remove, connect };
}
