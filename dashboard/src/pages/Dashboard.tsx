import { useQuery } from "@tanstack/react-query";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { apiMetrics } from "../api/metrics";

const formatFileSize = (size: number) => {
  if (size === null || size === 0) return "0 B";
  const i = Math.floor(Math.log(size) / Math.log(1024));
  return (
    (size / Math.pow(1024, i)).toFixed(2) +
    [" B", " KB", " MB", " GB", " TB"][i]
  );
};

const Dashboard = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["metrics"],
    queryFn: apiMetrics,
  });

  if (isLoading) return <p className="text-center text-gray-500">Loading...</p>;
  if (error)
    return <p className="text-center text-red-500">Error fetching metrics</p>;

  const metrics = data?.data?.metrics;

  // Prepare data for bar chart
  const fileTypes = Object.entries(metrics?.FilesByMimetype).map(
    ([type, size]) => ({
      type,
      size: size as number,
    })
  );

  return (
    <div className="p-6 space-y-9">
      {/* Title */}
      <h1 className="text-2xl font-bold text-gray-800 mb-5">📊 Dashboard Metrics</h1>

      {/* Metric Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card className="">
          <CardHeader>
            <CardTitle className="text-blue-700">📁 Total Files</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-xl font-semibold">{metrics?.TotalFileCount}</p>
          </CardContent>
        </Card>

        <Card className="">
          <CardHeader>
            <CardTitle className="text-green-700">🗑️ Deleted Files</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-xl font-semibold">{metrics?.TotalDeletedFiles}</p>
          </CardContent>
        </Card>

        <Card className="">
          <CardHeader>
            <CardTitle className="text-purple-700">💾 Total Storage</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-xl font-semibold">
              {formatFileSize(metrics?.TotalFileSize)}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* File Type Bar Chart */}
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-lg font-semibold mb-4">📦 File Distribution</h2>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={fileTypes}>
            <XAxis dataKey="type" tick={{ fontSize: 12 }} />
            <YAxis tickFormatter={formatFileSize} />
            <Tooltip formatter={(value) => formatFileSize(value as number)} />
            <Bar dataKey="size" fill="#4F46E5" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Last Updated Info */}
      <p className="text-gray-500 text-sm text-right">
        Last Updated: {new Date(metrics?.LastUpdated).toLocaleString()}
      </p>
    </div>
  );
};

export default Dashboard;
