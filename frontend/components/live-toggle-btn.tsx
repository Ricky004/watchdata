"use client"

import { Button } from "@/components/ui/button"

type LiveToggleButtonProps = {
    live: boolean
    toggleLive: () => void
}

export default function LiveToggleButton({ live, toggleLive }: LiveToggleButtonProps) {
    return (
        <Button
            onClick={toggleLive}
            className={`${live ? "bg-green-600 hover:bg-green-700" : "bg-red-600 hover:bg-red-700"} px-6 py-3 text-lg text-white`}
        >
            {live ? 'Live: on' : 'Live: off'}
        </Button>
    )
}
