import React, { useEffect, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { apiLogin } from "../api/auth";
import { useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader } from "../components/ui/card";

import lancer from "../assets/lancer.svg"
import { Label } from "../components/ui/label";
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";

import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup"
import { useForm } from "react-hook-form";

const schema = yup.object({
  email: yup.string().email("Must be a valid email").required("Required"),
  password: yup.string().required("Required")
})

const Login = () => {
  const [error, setError] = useState("");

  const { formState: { errors }, handleSubmit, register } = useForm({
    resolver: yupResolver(schema)
  })

  const navigate = useNavigate();
  const mutation = useMutation({
    mutationFn: apiLogin,
    onSuccess: (data) => {
      localStorage.setItem("token", data.data.token);
      navigate("/");
      alert("Login successful");
    },
    onError: (error) => {
      setError("Invalid email or password", error);
    },
  });

  useEffect(() => {
    console.log(errors)
  }, [errors])

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <Card className=" max-w-[550px] w-full">
        <CardHeader>
          <img src={lancer} className=" w-[65px] aspect-square" />
          <h2 className=" text-lg font-semibold">Welcome to Lancer</h2>
          <p>Login to admin dashboard to see your uploaded content</p>
        </CardHeader>
        <CardContent>
          <form className=" space-y-5" onSubmit={handleSubmit((data) => {
            setError("");
            mutation.mutate(data);
          })}>
            <div>
              <Label className=" block space-y-1">
                <span className=" block">Email</span>
                <Input placeholder="email@email.com"  {...register("email")} />
                {
                  errors["email"]?.message && <span className=" text-sm block text-red-500">{errors["email"]?.message}</span>
                }
              </Label>
            </div>
            <div>
              <Label className=" block space-y-1">
                <span className=" block">Password </span>
                <Input type="password" {...register("password")} />
                {
                  errors["password"]?.message && <span className=" text-sm block text-red-500">{errors["password"]?.message}</span>
                }
              </Label>
            </div>

            {
              error &&
              <div>
                <span className=" text-sm block text-red-500">{error}</span>
              </div>
            }

            <div className=" flex items-center justify-center">
              <Button>Login</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default Login;
