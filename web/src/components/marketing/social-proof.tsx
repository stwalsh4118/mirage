export function SocialProof() {
  const brands = ["Acme", "Globex", "Umbrella", "Soylent", "Initech", "Hooli"];
  return (
    <section className="py-12 section-soft border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-6 items-center opacity-70">
          {brands.map((b) => (
            <div key={b} className="text-center text-sm text-muted-foreground">{b}</div>
          ))}
        </div>
      </div>
    </section>
  );
}



