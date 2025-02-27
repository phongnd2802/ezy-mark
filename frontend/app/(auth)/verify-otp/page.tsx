"use client";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import VerifyOtpForm from "./otp-form";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function VerifyOtpPage() {
    const searchParams = useSearchParams();
    const token = searchParams.get("token")
    const [ttl, setTTL] = useState<number | null>(null);



    useEffect(() => {
        const getTTL = async () => {
            if (!token) return;
            try {
                const res = await fetch(`/api/v1/auth/verify-otp?token=${token}`);
                const data = await res.json();
                console.log(data);
                if (data.code === 20000) {

                    setTTL(parseInt(data.data.ttl));
                }
            } catch (error) {
                console.log(error);
            }
        }
        getTTL();
    }, [])


    useEffect(() => {
        if (ttl === null || ttl <= 0) return;

        const interval = setInterval(() => {
            setTTL(prevTTL => {
                if (prevTTL === null || prevTTL <= 1) {
                    clearInterval(interval);
                    return 0;
                }
                return prevTTL - 1;
            });
        }, 1000);

        return () => clearInterval(interval);
    }, [ttl]);


    
    return (
        <Card className="w-full max-w-md bg-white rounded-xl shadow-lg">
            <CardHeader>
                <CardTitle>Verify OTP - Complete your registration</CardTitle>
                <CardDescription className="text-justify">
                    Enter the OTP sent to your email to verify your identity and complete the registration process.
                </CardDescription>
            </CardHeader>
            <CardContent>

                <VerifyOtpForm token={token}/>

                <Separator />

                <div className="flex items-center justify-center space-x-1">
                    <p className="text-sm text-muted-foreground text-justify">
                        The OTP will expire in {" "}
                        <span className="text-red-500 font-semibold"> {ttl}</span> seconds.
                    </p>

                    <Button
                        variant="link"
                        className="px-0 py-0 text-blue-600 hover:text-blue-800 transition"
                        disabled={ttl !== null && ttl > 0}    
                    >
                        Resend
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}