import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function HowItWorks() {
  const steps = [
    { title: "Visualize your infrastructure", desc: "See all your Railway projects, environments, and services in a unified dashboard with real-time status." },
    { title: "Create with confidence", desc: "Use the guided environment wizard to provision new environments with the right configuration every time." },
    { title: "Monitor what matters", desc: "Track service health, deployment status, and get instant visibility into your entire Railway infrastructure." },
  ];
  return (
    <section className="py-24 section-soft border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight mb-4">Built for Railway power users</h2>
          <p className="text-xl text-muted-foreground">Everything you need to manage complex Railway infrastructure efficiently.</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {steps.map((s, i) => (
            <Card key={s.title} className="glass grain">
              <CardHeader>
                <CardTitle className="text-lg">{i + 1}. {s.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">{s.desc}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}



