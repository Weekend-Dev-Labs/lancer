import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "../components/ui/button";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "../components/ui/table";
import { toast } from "react-hot-toast";
import { apiDeleteUpload, apiUpload } from "../api/upload";

const UploadList = () => {
  const queryClient = useQueryClient();
  const [filters, setFilters] = useState({
    page: 1,
    count: 10,
    fileType: "",
    provider: "",
  });

  const { data, isLoading, error } = useQuery({
    queryKey: ["uploads", filters],
    queryFn: () => apiUpload(filters),
  });

  const deleteMutation = useMutation({
    mutationFn: apiDeleteUpload,
    onSuccess: () => {
      toast.success("File deleted successfully"),
        queryClient.invalidateQueries(["uploads"]);
    },
    onError: () => {
      toast.error("Failed to delete file");
    },
  });

  if (isLoading) return <p className="text-center text-gray-500">Loading...</p>;
  if (error)
    return <p className="text-center text-red-500">Error fetching uploads</p>;

  const uploads = data?.data?.uploads?.Files || [];

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold text-gray-800">📂 Uploaded Files</h1>

      {/* Filters */}
      <div className="flex space-x-4 my-4">
        <input
          type="text"
          placeholder="File Type"
          className="border p-2 rounded"
          value={filters.fileType}
          onChange={(e) =>
            setFilters((prev) => ({ ...prev, fileType: e.target.value }))
          }
        />
        <input
          type="text"
          placeholder="Provider"
          className="border p-2 rounded"
          value={filters.provider}
          onChange={(e) =>
            setFilters((prev) => ({ ...prev, provider: e.target.value }))
          }
        />
        <Button onClick={() => setFilters((prev) => ({ ...prev, page: 1 }))}>
          Apply
        </Button>
      </div>

      {/* File List */}
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>File Name</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Size</TableHead>
            <TableHead>Uploaded At</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {uploads.map((file) => (
            <TableRow key={file.ID}>
              <TableCell>{file.FileName}</TableCell>
              <TableCell>{file.FileType}</TableCell>
              <TableCell>{(file.FileSize / 1024).toFixed(2)} KB</TableCell>
              <TableCell>
                {new Date(file.UploadedAt).toLocaleString()}
              </TableCell>
              <TableCell>
                <Button
                  className="bg-red-500 text-white"
                  onClick={() => deleteMutation.mutate({ id: [file.ID] })}
                >
                  Delete
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      {/* Pagination */}
      <div className="flex justify-between items-center mt-4">
        <Button
          disabled={filters.page === 1}
          onClick={() =>
            setFilters((prev) => ({ ...prev, page: prev.page - 1 }))
          }
        >
          Previous
        </Button>
        <span>Page {filters.page}</span>
        <Button
          disabled={data?.data.uploads.Meta.TotalPage === filters.page}
          onClick={() =>
            setFilters((prev) => ({ ...prev, page: prev.page + 1 }))
          }
        >
          Next
        </Button>
      </div>
    </div>
  );
};

export default UploadList;
