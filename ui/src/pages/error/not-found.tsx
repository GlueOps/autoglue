import {useNavigate} from "react-router-dom";
import {Button} from "@/components/ui/button.tsx";

export const NotFoundPage = () => {
    const navigate = useNavigate();

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-background text-foreground">
            <h1 className="text-6xl font-bold mb-4">404</h1>
            <p className="text-2xl mb-8">Oops! Page not found</p>
            <Button onClick={() => navigate("/dashboard")}>
                Go back to Dashboard
            </Button>
        </div>
    );
};