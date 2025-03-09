import { AxiosResponse } from "axios";
import { axiosInstances } from "./axios";
import { getAuthHeaders } from "./auth";

interface ApiUpload {
  ID: string;
  FileName: string;
  FilePath: string;
  FileSize: number;
  FileType: string;
  UploadedBy: string;
  UploadedAt: string;
  IsDeleted: boolean;
  Checksum: string;
  Description: string;
  Provider: string;
  ProviderMetadata: string;
}

interface ApiUploadMeta {
  Page: number;
  Size: number;
  TotalPage: number;
  TotalCount: number;
}

interface ApiPayloadData {
  id: string[];
}
export const apiUpload = ({
  page,
  count,
  fileType,
  provider,
}: {
  page: number;
  count: number;
  fileType: string;
  provider: string;
}) => {
  const queryParams = [
    `page=${page}`,
    `size=${count}`,
    fileType ? `file_type=${fileType}` : "",
    provider ? `provider=${provider}` : "",
  ]
    .filter(Boolean)
    .join("&");

  return axiosInstances.get<
    any,
    AxiosResponse<{ uploads: { Files: ApiUpload[]; Meta: ApiUploadMeta } }>
  >(`upload?${queryParams}`, {
    headers: getAuthHeaders(),
  });
};

export const apiDeleteUpload = (data: ApiPayloadData) => {
  return axiosInstances.post(`/upload/delete`, data, {
    headers: getAuthHeaders(),
  });
};
