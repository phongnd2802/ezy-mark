import Image from "next/image";
import React from "react";




export default function AuthLayout(
    { children }: { children: React.ReactNode }
) {
    return (
        <div className="w-full flex">
            <div className="w-[50%] h-screen flex flex-col justify-center">
                <div className="flex px-16 mb-5 gap-x-2.5">
                    <Image
                        width={50}
                        height={50}
                        src="/logo.png"
                        alt="Logo"
                    />
                    <div className="h-[50px] flex items-center justify-center">
                        <h1 className="text-2xl font-bold">Daily Social</h1>
                    </div>
                </div>
                <div className="px-16">
                    <h2 className="text-4xl font-bold mb-2">Welcome to Daily Social</h2>
                    <p className="text-justify">
                        Daily Social is a platform that helps you stay connected with friends, keep up with the latest updates, and share your favorite moments.
                        Engage with your community, explore trending content, and create lasting memories—all in one place.
                    </p>
                </div>
            </div>
            <div className="w-[50%]">
                <div className="mx-10 h-screen flex items-center">
                    {children}
                </div>
            </div>
        </div>
    )
}