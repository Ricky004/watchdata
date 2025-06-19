"use client"

import { Button } from "@/components/ui/button"
import { Play, Pause } from 'lucide-react';

type LiveToggleButtonProps = {
    live: boolean
    toggleLive: () => void
}

export default function LiveToggleButton({ live, toggleLive }: LiveToggleButtonProps) {
    return (
        <Button
            onClick={toggleLive}
            className={`${live ? "bg-red-600 hover:bg-red-700" : "bg-gray-600 hover:bg-gray-700"} px-6 py-3 text-base text-white`}
        >
            {live ? (
                <>
                  <Pause />
                   Live
                </>
            ) : (
                <>
                    <Play />
                    Live
                </>
            )}
        </Button>
    )
}
