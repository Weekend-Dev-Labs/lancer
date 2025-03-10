import { AxiosResponse } from "axios";
import { axiosInstances } from "./axios";
import { getAuthHeaders } from "./auth";

export interface APISettings {
  allowedOrigins: string[];
  authWebhook: string;
  awsBucket: string;
  awsRegion: string;
  database: string;
  databaseName: string;
  eventsWebhook: string;
  isAuthenticationEnabled: boolean;
  isAwsEnabled: boolean;
  isRedis: boolean;
  port: string;
  redisServer: string;
  storePath: string;
  tempPath: string;
  webhookSecret: string;
}

export const apiGetSettings = () => {
  return axiosInstances.get<any, AxiosResponse<APISettings>>("/settings", {
    headers: getAuthHeaders(),
  });
};
