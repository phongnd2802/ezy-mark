"use client";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { authService } from "@/services/authService";
import { ERROR_CODES } from "@/utils/errorCode";
import { ERROR_MESSAGES } from "@/utils/errorMessages";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";

import { useForm } from "react-hook-form";
import { z } from "zod";


const signUpSchema = z.object({
    email: z.string().email({ message: "Invalid email address" }),
    password: z.string().min(6, { message: "Password must be at 6 least characters" }),
    confirmPassword: z.string()
})

export default function SignUpForm() {
    const router = useRouter();
    const form = useForm<z.infer<typeof signUpSchema>>({
        resolver: zodResolver(signUpSchema),
        defaultValues: {
            email: "",
            password: "",
            confirmPassword: "",
        },
    })


    const onSubmit = async (values: z.infer<typeof signUpSchema>) => {
        if (values.password !== values.confirmPassword) {
            form.setError("confirmPassword", { message: "Passwords do not match" });
            form.setValue("confirmPassword", "");
            return;
        }
        try {
            
            const result = await authService.signUp({
                email: values.email,
                password: values.password,
            })
            if (result.code === 20000) {
                const { token } = result.data;
                router.push(`/verify-otp?token=${token}`)
            } else if (result.code === ERROR_CODES.EMAIL_ALREADY_EXISTS) {
                form.setError("email", { message: ERROR_MESSAGES[result.code] });

                form.setValue("password", "");
                form.setValue("confirmPassword", "");
            } else if (result.code === ERROR_CODES.PENDING_VERIFICATION) {
                const { token } = result.data;
                router.push(`/verify-otp?token=${token}`)
            }
        } catch (error) {
            console.log("Error:", error);
        }
    }

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-y-2.5">
                <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Email</FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>

                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Password</FormLabel>
                            <FormControl>
                                <Input type="password" {...field} />
                            </FormControl>

                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="confirmPassword"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Confirm Password</FormLabel>
                            <FormControl>
                                <Input type="password" {...field} />
                            </FormControl>

                            <FormMessage />
                        </FormItem>
                    )}
                />

                <Button type="submit" className="w-full my-2.5">
                    Continue
                </Button>
            </form>
        </Form>
    )
}