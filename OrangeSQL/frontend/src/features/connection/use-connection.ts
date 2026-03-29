import { ref, onMounted } from "vue";
import {
  fetchProfiles,
  createProfile,
  updateProfile,
  deleteProfile,
  connectProfile,
} from "../../shared/api";
import type { ConnectionProfile, CreateProfileRequest, UpdateProfileRequest } from "./types";

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
      profiles.value = res.profiles as ConnectionProfile[];
      activeId.value = res.activeId;
    } catch {
      // エラー時は何もしない
    } finally {
      loading.value = false;
    }
  }

  async function create(req: CreateProfileRequest): Promise<void> {
    await createProfile(req);
    await refresh();
  }

  async function update(id: string, req: UpdateProfileRequest): Promise<void> {
    await updateProfile(id, req);
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
