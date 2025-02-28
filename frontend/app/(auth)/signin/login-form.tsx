"use client";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { authService } from "@/services/authService";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { ERROR_CODES } from "@/utils/errorCode";
import { ERROR_MESSAGES } from "@/utils/errorMessages";

const loginSchema = z.object({
    email: z.string().email({ message: "Invalid email address" }),
    password: z.string().min(6, { message: "Password must be at least 6 characters" }),
})



export default function LoginForm() {
    const router = useRouter();

    const form = useForm<z.infer<typeof loginSchema>>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    })

    const onSubmit = async (values: z.infer<typeof loginSchema>) => {
        try {
            const result = await authService.logIn({
                email: values.email,
                password: values.password,
            })
            if(result.code === 20000) {
                Cookies.set("access-token", result.data.access_token);
                Cookies.set("refresh-token", result.data.refresh_token);

                router.replace("/")
            } else if (result.code === ERROR_CODES.AUTHENTICATION_FAILED) {
                form.setError("password", {message: ERROR_MESSAGES[result.code]});
                form.setValue("password", "");
            } else if (result.code === ERROR_CODES.ACCOUNT_NOT_VERIFIED) {
                form.setError("password", {message: ERROR_MESSAGES[result.code]})
                form.setValue("password", "");
            }

        } catch (error) {
            console.log(error);
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

                <Button type="submit" className="w-full my-2.5">
                    Continue
                </Button>
            </form>
        </Form>
    )
}