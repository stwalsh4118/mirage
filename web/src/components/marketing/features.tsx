import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function Features() {
  const features = [
    { title: "Project Dashboard", description: "View all your Railway projects, environments, and services in one unified interface.", icon: "📊" },
    { title: "Environment Wizard", description: "Create new environments with a guided wizard that simplifies complex Railway configuration.", icon: "🪄" },
    { title: "Real-time Updates", description: "Monitor project and service status with live updates powered by Railway's API.", icon: "⚡" },
    { title: "Secure by Default", description: "Built-in authentication with Clerk and secure management of Railway API tokens.", icon: "🔐" },
    { title: "Railway-Native", description: "Deep integration with Railway's GraphQL API for seamless infrastructure management.", icon: "🚄" },
    { title: "Multi-Project Support", description: "Manage multiple Railway projects and their environments from a single dashboard.", icon: "🎯" },
  ];

  return (
    <section className="py-24 relative section-soft-alt border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight mb-4">Powerful Railway management made simple</h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">Everything you need to manage your Railway infrastructure efficiently in one beautiful interface.</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 gap-6">
          {features.map((feature, index) => {
            const getColSpan = (i: number) => (i === 0 ? "lg:col-span-5" : i === 1 ? "lg:col-span-4" : i === 2 ? "lg:col-span-3" : "lg:col-span-4");
            const getGlass = (i: number) => (i % 3 === 0 ? "bg-card/80 backdrop-blur-xl border-border/60 shadow-lg" : i % 2 === 0 ? "bg-card/65 backdrop-blur-md border-border/50" : "bg-card/40 backdrop-blur-sm border-border/30");
            return (
              <Card key={index} className={`${getGlass(index)} ${getColSpan(index)} grain border hover:scale-[1.02] hover:-translate-y-2 transition-all duration-300 shadow-lg hover:shadow-2xl sheen`}>
                <CardHeader>
                  <div className="text-2xl mb-2">{feature.icon}</div>
                  <CardTitle className="text-lg">{feature.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground leading-relaxed">{feature.description}</p>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </div>
    </section>
  );
}



