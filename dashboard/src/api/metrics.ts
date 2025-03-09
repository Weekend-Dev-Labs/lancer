import { AxiosResponse } from "axios";
import { axiosInstances } from "./axios";
import { getAuthHeaders } from "./auth";

interface ApiMatricsResponse {
  FilesByMimetype: Record<string, number>;
  ID: string;
  LargestFileSize: number;
  LastUpdated: string;
  SmallestFileSize: null;
  TotalDeletedFiles: number;
  TotalFileCount: number;
  TotalFileSize: number;
}

export const apiMetrics = () => {
  return axiosInstances.get<
    any,
    AxiosResponse<{ metrics: ApiMatricsResponse }>
  >("/metrics", { headers: getAuthHeaders() });
};
