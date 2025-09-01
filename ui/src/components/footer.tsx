import { FaGithub } from "react-icons/fa";

export function Footer() {
    return (
        <footer className="border-t">
            <div className="container flex flex-col items-center justify-between gap-4 py-10 md:h-24 md:flex-row md:py-0">
                <div className="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
                    <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
                        Built for{" "}
                        <a
                            href="https://www.glueops.dev/"
                            target="_blank"
                            rel="noreferrer"
                            className="font-medium underline underline-offset-4"
                        >
                            GlueOps
                        </a>
                        . The source code is available on{" "}
                        <a
                            href="https://github.com/GlueOps/autoglue"
                            target="_blank"
                            rel="noreferrer"
                            className="font-medium underline underline-offset-4"
                        >
                            GitHub
                        </a>
                        .
                    </p>
                </div>
                <div className="flex items-center space-x-4">
                    <a href="https://github.com/GlueOps/autoglue" target="_blank" rel="noreferrer">
                        <FaGithub className="h-5 w-5" />
                    </a>
                </div>
            </div>
        </footer>
    );
}
