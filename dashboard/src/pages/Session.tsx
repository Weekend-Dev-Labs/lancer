import React from "react";
import { useQuery } from "@tanstack/react-query";
import { apiSession } from "../api/session";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Skeleton } from "../components/ui/skeleton";

const Session = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["sessions"],
    queryFn: apiSession,
  });

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {[...Array(6)].map((_, index) => (
          <Skeleton key={index} className="h-32 w-full rounded-lg" />
        ))}
      </div>
    );
  }

  if (error) {
    return <p className="text-red-500 text-center">Failed to load sessions.</p>;
  }

  return (
    <div className="p-6">
      <h2 className="text-2xl font-semibold mb-4">Active Sessions</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {data?.data?.sessions?.map((session) => (
          <Card key={session.id} className="shadow-md hover:shadow-lg transition">
            <CardHeader>
              <CardTitle className="text-lg">{session.file_name}</CardTitle>
              <p className="text-sm text-gray-500">Provider: {session.provider}</p>
            </CardHeader>
            <CardContent>
              <p className="text-gray-700">
                <span className="font-semibold">Size:</span> {session.file_size} bytes
              </p>
              <p className="text-gray-700">
                <span className="font-semibold">Chunks:</span> {session.current_chunk}/{session.max_chunk}
              </p>
              <p className="text-gray-700">
                <span className="font-semibold">Owner:</span> {session.owner_id}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Session;
