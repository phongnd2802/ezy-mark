"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";

export default function VerifyOtpTimer({ ttl: initialTTL }: { ttl: number }) {
    const [ttl, setTTL] = useState<number | null>(initialTTL);

    useEffect(() => {
        if (ttl === null || ttl <= 0) return;

        const interval = setInterval(() => {
            setTTL((prevTTL) => {
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
        <div className="flex items-center justify-center space-x-1">
            <p className="text-sm text-muted-foreground text-justify">
                The OTP will expire in <span className="text-red-500 font-semibold">{ttl}</span> seconds.
            </p>
            <Button
                variant="link"
                className="px-0 py-0 text-blue-600 hover:text-blue-800 transition"
                disabled={ttl !== null && ttl > 0}
            >
                Resend
            </Button>
        </div>
    );
}
