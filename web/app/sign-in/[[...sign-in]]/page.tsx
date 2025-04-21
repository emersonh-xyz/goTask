import { SignIn } from '@clerk/nextjs'

export default function Page() {
    return (
        <div className="h-screen flex items-center justify-center bg-base-300 gap-8">
            <h1 className="text-5xl font-bold">
                go<span className="text-primary">Task</span>.
            </h1>
            <SignIn />
        </div>
    )
}