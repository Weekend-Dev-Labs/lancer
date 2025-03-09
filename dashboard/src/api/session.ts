import { AxiosResponse } from "axios";
import { axiosInstances } from "./axios";
import { getAuthHeaders } from "./auth";

interface ApiSessionResponse {
  id: string;
  file_size: number;
  chunk_size: number;
  max_chunk: number;
  file_name: string;
  temp_path: string;
  owner_id: string;
  current_chunk: number;
  provider: string;
}

export const apiSession = () => {
  return axiosInstances.get<
    any,
    AxiosResponse<{ sessions: ApiSessionResponse[] }>
  >("/sessions", { headers: getAuthHeaders() });
};
