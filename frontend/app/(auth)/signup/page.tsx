import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import Link from "next/link";
import SignUpForm from "./signup-form";
import { Metadata } from "next";


export const metadata: Metadata = {
    title: "Sign up",
  }

export default function SignUpPage() {

    return (
        <Card className="w-full max-w-md bg-white rounded-xl shadow-lg">
            <CardHeader>
                <CardTitle>Sign up to continue</CardTitle>
                <CardDescription>Use your email or another service to continue</CardDescription>
            </CardHeader>
            <CardContent>
                <SignUpForm />

                <div className="flex flex-col">
                    <div className="w-full text-center">
                        <span>
                            Already have an account? {" "}
                            <Link
                                href="/signin"
                                className="text-blue-500 hover:underline hover:text-blue-600 transition duration-200"
                            >
                                Log in
                            </Link>
                        </span>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}