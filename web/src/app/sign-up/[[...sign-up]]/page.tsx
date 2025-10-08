import { SignUp } from "@clerk/nextjs";

export default function SignUpPage() {
	return (
		<div className="flex min-h-screen items-center justify-center bg-background px-4 py-12">
			<div className="w-full max-w-md space-y-8">
				<div className="text-center">
					<h1 className="text-4xl font-bold tracking-tight">Create your account</h1>
					<p className="mt-2 text-muted-foreground">
						Start provisioning ephemeral environments in seconds
					</p>
				</div>

				<div className="flex justify-center">
					<SignUp
						appearance={{
							elements: {
								rootBox: "w-full",
								card: "shadow-lg border border-border rounded-lg",
								headerTitle: "text-2xl font-semibold",
								headerSubtitle: "text-muted-foreground",
								socialButtonsBlockButton:
									"border border-border hover:bg-accent hover:text-accent-foreground transition-colors",
								formButtonPrimary:
									"bg-primary hover:bg-primary/90 text-primary-foreground transition-colors",
								formFieldInput:
									"border-border focus:ring-2 focus:ring-primary/20 focus:border-primary",
								footerActionLink:
									"text-primary hover:text-primary/80 transition-colors",
								dividerLine: "bg-border",
								dividerText: "text-muted-foreground",
								formFieldLabel: "text-foreground",
								identityPreviewText: "text-foreground",
								formResendCodeLink:
									"text-primary hover:text-primary/80 transition-colors",
								otpCodeFieldInput:
									"border-border focus:ring-2 focus:ring-primary/20 focus:border-primary",
							},
						}}
					/>
				</div>
			</div>
		</div>
	);
}

