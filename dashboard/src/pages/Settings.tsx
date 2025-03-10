import { useMemo, useState } from "react";
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
import { useNavigate } from "react-router-dom";
import { apiGetSettings } from "../api/settings";
import { formatText } from "../lib/utils";

const Settings = () => {
    const queryClient = useQueryClient();

    const { data, isLoading, error } = useQuery({
        queryKey: ["settings"],
        queryFn: async () => (await apiGetSettings()).data,
    });

    const settings = useMemo(() => {
        if (data) {
            const value = Object.keys(data).map((key) => {
                const keyValue = (data as any)[key];
                if (Array.isArray(keyValue)) {
                    return { setting: formatText(key), value: keyValue.join(" , ") }
                }

                return { setting: formatText(key), value: keyValue }
            })

            return value;
        }
    }, [data])

    const navigate = useNavigate();

    if (isLoading) return <p className="text-center text-gray-500">Loading...</p>;
    if (error) {
        localStorage.clear();
        navigate("/login")
        return <p className="text-center text-red-500">Error fetching lancer settings</p>;
    }

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold text-gray-800">⚙️ Uploaded Files</h1>

            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Setting</TableHead>
                        <TableHead>Value</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {
                        settings?.map((val) => (

                            <TableRow key={val.setting}>
                                <TableCell>{val.setting}</TableCell>
                                <TableCell>{val.value}</TableCell>
                            </TableRow>
                        ))
                    }
                </TableBody>
            </Table>
        </div>
    );
};

export default Settings;
