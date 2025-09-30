import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

type Kpi = { title: string; value: number; delta?: number };

export function KpiStrip({ items }: { items: Kpi[] }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {items.map((k) => (
        <Card key={k.title} className="glass grain">
          <CardHeader className="pb-2">
            <CardTitle className="text-xs text-muted-foreground font-medium flex items-center justify-between">
              {k.title}
              {typeof k.delta === "number" && (
                <span className="text-[10px] rounded-full px-1.5 py-0.5 bg-muted/60">{k.delta >= 0 ? `+${k.delta}` : k.delta}</span>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="text-2xl font-semibold">{k.value}</div>
            <div className="mt-2 h-6 rounded bg-muted/50" />
          </CardContent>
        </Card>
      ))}
    </div>
  );
}






