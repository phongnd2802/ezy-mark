"use client";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { InputOTP, InputOTPSlot } from "@/components/ui/input-otp";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS } from "input-otp";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const otpSchema = z.object({
    otp: z.string().length(6, { message: "Your OTP must be exactly 6 digits" }),
});

export default function VerifyOtpForm(props: { token: string | null }) {
    const router = useRouter();

    const form = useForm<z.infer<typeof otpSchema>>({
        resolver: zodResolver(otpSchema),
        defaultValues: { otp: "" },
    });

    const [loading, setLoading] = useState(false);

    const onSubmit = async (data: z.infer<typeof otpSchema>) => {
        setLoading(true);
        try {
            const res = await fetch("/api/v1/auth/verify-otp", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    token: props.token,
                    otp: data.otp,
                }),
            })

            const result = await res.json();
            if (result.code === 20000) {

                router.replace("/signin");
            }
        } catch (error) {
            console.log(error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="w-full flex flex-col items-center gap-2.5">
                <FormField
                    control={form.control}
                    name="otp"
                    render={({ field }) => (
                        <FormItem className="text-center">
                            <FormControl>
                                <InputOTP
                                    maxLength={6} {...field}
                                    className="flex gap-2"
                                    pattern={REGEXP_ONLY_DIGITS}
                                >
                                    {[...Array(6)].map((_, index) => (
                                        <InputOTPSlot
                                            key={index}
                                            index={index}
                                            className="w-12 h-12 text-xl font-bold text-center border border-gray-300 rounded-md outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
                                        />
                                    ))}
                                </InputOTP>
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <Button type="submit" className="w-full my-2.5" disabled={loading}>
                    {loading ? "Verifying..." : "Continue"}

                </Button>
            </form>
        </Form>
    );
}
