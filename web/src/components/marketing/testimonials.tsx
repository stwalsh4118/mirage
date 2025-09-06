export function Testimonials() {
  const items = [
    { name: "Alex P.", role: "Platform Eng.", quote: "Mirage cut our env setup from hours to minutes." },
    { name: "Sam R.", role: "Staff Eng.", quote: "The glass UI and status updates are fantastic." },
  ];
  return (
    <section className="py-24 section-soft border-t border-border/40">
      <div className="container mx-auto px-4 grid gap-6 md:grid-cols-2">
        {items.map((t) => (
          <figure key={t.name} className="glass grain rounded-xl p-8">
            <blockquote className="text-lg text-foreground/90">“{t.quote}”</blockquote>
            <figcaption className="mt-4 text-sm text-muted-foreground">{t.name} — {t.role}</figcaption>
          </figure>
        ))}
      </div>
    </section>
  );
}



