export function Testimonials() {
  const items = [
    { name: "Alex P.", role: "Platform Engineer", quote: "Finally, a clean interface for managing our Railway infrastructure. The project overview saves us hours every week." },
    { name: "Sam R.", role: "Staff Engineer", quote: "Mirage makes it easy to see all our Railway projects and environments in one place. The real-time status updates are invaluable." },
  ];
  return (
    <section className="py-24 section-soft border-t border-border/40">
      <div className="container mx-auto px-4 grid gap-6 md:grid-cols-2">
        {items.map((t) => (
          <figure key={t.name} className="glass grain rounded-xl p-8">
            <blockquote className="text-lg text-foreground/90">&ldquo;{t.quote}&rdquo;</blockquote>
            <figcaption className="mt-4 text-sm text-muted-foreground">{t.name} — {t.role}</figcaption>
          </figure>
        ))}
      </div>
    </section>
  );
}



