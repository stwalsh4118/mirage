import { Header } from "@/components/marketing/header";
import { Hero } from "@/components/marketing/hero";
import { SocialProof } from "@/components/marketing/social-proof";
import { Features } from "@/components/marketing/features";
import { HowItWorks } from "@/components/marketing/how-it-works";
import { EnvironmentCards } from "@/components/marketing/environment-cards";
import { Testimonials } from "@/components/marketing/testimonials";
import { FinalCTA } from "@/components/marketing/final-cta";
import { Footer } from "@/components/marketing/footer";

export default function HomePage() {
  return (
    <div className="min-h-screen relative sandstorm-bg">
      <Header />
      <main className="relative z-10">
        <Hero />
        <SocialProof />
        <Features />
        <HowItWorks />
        <EnvironmentCards />
        <Testimonials />
        <FinalCTA />
      </main>
      <Footer />
    </div>
  );
}
