import axios from "axios";

export const axiosInstances = axios.create({
  baseURL: "https://add4-49-47-8-82.ngrok-free.app/api",
});
