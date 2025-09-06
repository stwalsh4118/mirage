import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function HowItWorks() {
  const steps = [
    { title: "Connect repo", desc: "Point Mirage at your monorepo or service repo." },
    { title: "Pick a template", desc: "Choose Dev or Prod presets, tweak as needed." },
    { title: "Create environment", desc: "Provision infra and deploy automatically." },
  ];
  return (
    <section className="py-24 section-soft border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight mb-4">How it works</h2>
          <p className="text-xl text-muted-foreground">From repo to running environment in minutes.</p>
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



