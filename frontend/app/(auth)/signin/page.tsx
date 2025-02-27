import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import LoginForm from "./login-form";
import { Button } from "@/components/ui/button";
import { FcGoogle } from "react-icons/fc";
import { FaGithub } from "react-icons/fa";
import Link from "next/link";
import { Metadata } from "next";


export const metadata: Metadata = {
  title: "Login",
}

export default function SignIn() {
  return (
    <Card className="w-full max-w-md bg-white rounded-xl shadow-lg">
      <CardHeader>
        <CardTitle>Login to continue</CardTitle>
        <CardDescription>Use your email or another service to continue</CardDescription>
      </CardHeader>
      <CardContent>
        <LoginForm />
        <div className="relative flex items-center w-full mb-2.5">
          <div className="flex-grow border-t border-gray-300"></div>
          <span className="px-4 text-gray-500 text-sm font-medium">OR</span>
          <div className="flex-grow border-t border-gray-300"></div>
        </div>

        <div className="flex flex-col gap-y-2.5">
          <Button variant="outline" className="flex items-center gap-2">
            <FcGoogle size={16} />
            Continue with Google
          </Button>

          <Button variant="outline" className="flex items-center gap-2">
            <FaGithub size={16} />
            Continue with Github
          </Button>

          <div className="w-full text-center">
            <Link
              href="/forgot-password"
              className="text-gray-400 hover:text-gray-500 transition duration-200"
            >
              Forgot password?
            </Link>
          </div>

          <div className="w-full text-center">
            <span>
              Don&apos;t have an account? {" "}
              <Link
                href="signup"
                className="text-blue-500 hover:underline hover:text-blue-600 transition duration-200"
              >
                Sign up
              </Link>
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}