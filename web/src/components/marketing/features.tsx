import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function Features() {
  const features = [
    { title: "One-click environments", description: "Deploy complete environments with a single click. No complex configuration required.", icon: "⚡" },
    { title: "Dev/Prod templates", description: "Pre-configured templates that mirror your production setup perfectly.", icon: "🎯" },
    { title: "Real-time status", description: "Monitor environment health and deployment progress in real-time.", icon: "📊" },
    { title: "Secure tokens", description: "Automatic token management and secure environment variable handling.", icon: "🔐" },
    { title: "Railway-native", description: "Built specifically for Railway with deep platform integration.", icon: "🚄" },
    { title: "Monorepo-aware", description: "Smart detection and deployment of monorepo services and dependencies.", icon: "📦" },
  ];

  return (
    <section className="py-24 relative section-soft-alt border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight mb-4">Everything you need for perfect environments</h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">From development to production, Mirage handles the complexity so you can focus on building.</p>
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



