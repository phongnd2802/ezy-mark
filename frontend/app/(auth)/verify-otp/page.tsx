import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import VerifyOtpForm from "./otp-form";
import { Separator } from "@/components/ui/separator";
import { redirect } from "next/navigation";
import { authService } from "@/services/authService";
import VerifyOtpTimer from "./verify-otp-timer";

export default async function VerifyOtpPage({ searchParams }: { searchParams: { token?: string } }) {
    const { token } = await searchParams;
    if (!token) {
        redirect("/signin")
    }
    try {
        const result = await authService.getTTLOtp(token);
        if (result.code !== 20000) {
            redirect("/signin");
        }
        const ttl = parseInt(result.data.ttl);

        return (
            <Card className="w-full max-w-md bg-white rounded-xl shadow-lg">
                <CardHeader>
                    <CardTitle>Verify OTP - Complete your registration</CardTitle>
                    <CardDescription className="text-justify">
                        Enter the OTP sent to your email to verify your identity and complete the registration process.
                    </CardDescription>
                </CardHeader>
                <CardContent>

                    <VerifyOtpForm token={token} />

                    <Separator />

                    <VerifyOtpTimer ttl={ttl} />
                </CardContent>
            </Card>
        )
    } catch (error) {
        console.log(error);
        redirect("/signin");
    }
}